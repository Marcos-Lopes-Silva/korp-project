import { Routes } from '@angular/router';
import { StockComponent } from './features/stock/stock.component';
import { App } from './app';
import { HomeComponent } from './features/home/home.component';
import { BillingComponent } from './features/billing/billing.component';

export const routes: Routes = [
  {path: '', component: HomeComponent},
  {path: 'estoque', component: StockComponent},
  {path: 'faturamento', component: BillingComponent}
];
