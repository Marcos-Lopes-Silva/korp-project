import { Component, inject, signal, computed, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatChipsModule } from '@angular/material/chips';
import { MatSnackBar } from '@angular/material/snack-bar';

import { StockService } from '../../core/services/stock.service';
import { Product } from '../../core/model/product';
import { Invoice } from '../../core/model/invoice';
import { InvoiceService } from '../../core/services/invoice.service';

@Component({
  selector: 'app-billing',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    MatIconModule,
    MatTableModule,
    MatChipsModule,
  ],
  templateUrl: './billing.component.html',
})
export class BillingComponent implements OnInit {
  private invoiceService = inject(InvoiceService);
  private stockService = inject(StockService);
  private snackBar = inject(MatSnackBar);

  invoices = signal<Invoice[]>([]);
  products = signal<Product[]>([]);
  searchTerm = signal('');

  selectedInvoice = signal<Invoice | null>(null);
  newItem = { productId: '', quantity: 1 };

  formModel = { name: '', address: '' };

  filteredInvoices = computed(() => {
    const term = this.searchTerm().toLowerCase();
    return this.invoices().filter(
      (i) =>
        i.name.toLowerCase().includes(term) ||
        String(i.seqCode ?? '').includes(term)
    );
  });

  private reservedInCurrentInvoice = computed(() => {
    const reserved = new Map<string, number>();
    const invoice = this.selectedInvoice();
    if (!invoice?.items) return reserved;

    for (const item of invoice.items) {
      reserved.set(item.productId, (reserved.get(item.productId) ?? 0) + item.quantity);
    }
    return reserved;
  });

  availableStock(): number | null {
    if (!this.newItem.productId) return null;
    const product = this.products().find((p) => p.id === this.newItem.productId);
    if (!product) return null;
    const alreadyReserved = this.reservedInCurrentInvoice().get(product.id!) ?? 0;
    return product.quantity - alreadyReserved;
  }

  displayedColumns: string[] = ['seqCode', 'name', 'status', 'items', 'total', 'actions'];

  ngOnInit() {
    this.loadInvoices();
    this.loadProducts();
  }

  loadInvoices() {
    this.invoiceService.list().subscribe({
      next: (data) => this.invoices.set(data),
      error: () => this.notify('Não foi possível carregar as notas.', 'error'),
    });
  }

  loadProducts() {
    this.stockService.list().subscribe({
      next: (data) => this.products.set(data),
      error: () => this.notify('Não foi possível carregar os produtos.', 'error'),
    });
  }

  productName(productId: string): string {
    return this.products().find((p) => p.id === productId)?.name ?? productId;
  }

  createInvoice() {
    if (!this.formModel.name || !this.formModel.address) return;

    this.invoiceService.create(this.formModel).subscribe({
      next: () => {
        this.formModel = { name: '', address: '' };
        this.loadInvoices();
        this.notify('Nota criada com sucesso.', 'success');
      },
      error: () => this.notify('Erro ao criar a nota.', 'error'),
    });
  }

  manageItems(invoice: Invoice) {
    this.selectedInvoice.set(invoice);
    this.newItem = { productId: '', quantity: 1 };
  }

  closePanel() {
    this.selectedInvoice.set(null);
  }

  addItem() {
    const invoice = this.selectedInvoice();
    if (!invoice?.id || !this.newItem.productId) return;

    const available = this.availableStock();

    if (available === null) {
      this.notify('Selecione um produto.', 'error');
      return;
    }

    if (this.newItem.quantity <= 0) {
      this.notify('Informe uma quantidade válida.', 'error');
      return;
    }

    if (this.newItem.quantity > available) {
      this.notify(`Saldo insuficiente. Disponível: ${available} unidade(s).`, 'error');
      return;
    }

    this.invoiceService
      .addItem(invoice.id, this.newItem.productId, this.newItem.quantity)
      .subscribe({
        next: (updated) => {
          this.selectedInvoice.set(updated);
          this.newItem = { productId: '', quantity: 1 };
          this.loadInvoices();
          this.notify('Item adicionado.', 'success');
        },
        error: () => this.notify('Erro ao adicionar item — verifique o estoque.', 'error'),
      });
  }

  printInvoice(invoice: Invoice) {
    if (!invoice.id) return;

    this.invoiceService.print(invoice.id).subscribe({
      next: (blob: Blob) => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `invoice_${invoice.id}.pdf`;
        document.body.appendChild(a);
        a.click();

        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);
        this.closePanel();
        this.loadInvoices();
        this.notify('Nota impressa e fechada.', 'success');
      },
      error: (error) => {
        console.error(error)
        this.notify('Erro ao imprimir a nota. Verifique o estoque.', 'error')
      }
    });
  }

  private notify(message: string, kind: 'success' | 'error') {
    this.snackBar.open(message, 'Fechar', {
      duration: 3500,
      panelClass: kind === 'success' ? 'toast-success' : 'toast-error',
      horizontalPosition: 'right',
      verticalPosition: 'bottom',
    });
  }
}
