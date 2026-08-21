import { Component, inject, signal, computed } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { StockService } from '../../core/services/stock.service';
import { Product } from '../../core/model/product';

import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatChipsModule } from '@angular/material/chips';
import { CommonModule } from '@angular/common';


@Component({
  selector: 'app-stock',
  standalone: true,
  imports: [
    FormsModule,
    CommonModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatTableModule,
    MatChipsModule,
  ],
  templateUrl: './stock.component.html'
})
export class StockComponent {
  private stockService = inject(StockService);

  products = signal<Product[]>([]);
  searchTerm = signal('');
  editingProduct = signal<Product | null>(null);

  formModel: Product = this.emptyProduct();

  filteredProducts = computed(() => {
    const term = this.searchTerm().toLowerCase();
    return this.products().filter(p =>
      p.name.toLowerCase().includes(term) || p.sku.toLowerCase().includes(term)
    );
  });

  ngOnInit() {
    this.load();
  }

  load() {
    this.stockService.list().subscribe(data => this.products.set(data));
  }

  emptyProduct(): Product {
    return { name: '', sku: '', quantity: 0, price: 0 };
  }

  startEdit(product: Product) {
    this.editingProduct.set(product);
    this.formModel = { ...product };
  }

  cancelEdit() {
    this.editingProduct.set(null);
    this.formModel = this.emptyProduct();
  }

  save() {
    const editing = this.editingProduct();
    if (editing && editing.id) {
      this.stockService.update(editing.id, this.formModel).subscribe(() => {
        this.cancelEdit();
        this.load();
      });
    } else {
      this.stockService.create(this.formModel).subscribe(() => {
        this.cancelEdit();
        this.load();
      });
    }
  }

  remove(product: Product) {
    if (!product.id) return;
    if (!confirm(`Remover "${product.name}"?`)) return;
    this.stockService.delete(product.id).subscribe(() => this.load());
  }

  displayedColumns: string[] = ['name', 'sku', 'quantity', 'price', 'actions'];
}
