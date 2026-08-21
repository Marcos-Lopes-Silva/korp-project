import { Component, inject, signal, computed, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';

import { StockService } from '../../core/services/stock.service';
import { InvoiceService } from '../../core/services/invoice.service';
import { Product } from '../../core/model/product';
import { Invoice } from '../../core/model/invoice';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './home.component.html',
})
export class HomeComponent implements OnInit {
  private stockService = inject(StockService);
  private invoiceService = inject(InvoiceService);

  userName = 'Marcos';

  today = new Date().toLocaleDateString('pt-BR', {
    weekday: 'long',
    day: '2-digit',
    month: 'long',
  });

  products = signal<Product[]>([]);
  invoices = signal<Invoice[]>([]);

  produtosCadastrados = computed(() => this.products().length);

  notasAbertas = computed(
    () => this.invoices().filter((i) => i.status === 'Aberta').length
  );

  notasFechadas = computed(
    () => this.invoices().filter((i) => i.status === 'Fechada').length
  );

  faturadoNoMes = computed(() =>
    this.invoices()
      .filter((i) => i.status === 'Fechada')
      .reduce((total, i) => total + (i.totalPrice ?? 0), 0)
  );

  ngOnInit() {
    this.stockService.list().subscribe((data) => this.products.set(data));
    this.invoiceService.list().subscribe((data) => this.invoices.set(data));
  }
}
