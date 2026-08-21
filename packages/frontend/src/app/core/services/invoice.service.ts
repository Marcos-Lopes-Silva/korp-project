import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, map } from 'rxjs';
import { Invoice } from '../model/invoice';
import { environment } from '../../environments/environment';


@Injectable({ providedIn: 'root' })
export class InvoiceService {
  private baseUrl = `${environment.apiUrl}/invoices`;

  constructor(private http: HttpClient) { }

  list(): Observable<Invoice[]> {
    return this.http.get<Invoice[]>(this.baseUrl)
  }

  getById(id: string): Observable<Invoice> {
    return this.http.get<Invoice>(`${this.baseUrl}/${id}`)
  }

  create(invoice: Pick<Invoice, 'name' | 'address'>): Observable<Invoice> {
    return this.http.post<Invoice>(this.baseUrl, invoice)
  }

  addItem(invoiceId: string, productId: string, quantity: number): Observable<Invoice> {
    return this.http
      .post<Invoice>(`${this.baseUrl}/${invoiceId}/items`, { product_id: productId, quantity })
  }

  print(invoiceId: string): Observable<Blob> {
    return this.http.post<Blob>(`${this.baseUrl}/${invoiceId}/print`, {}, { responseType: 'blob' })
  }
}
