export interface InvoiceItem {
  id?: string;
  invoiceId?: string;
  productId: string;
  quantity: number;
}

export interface Invoice {
  id?: string;
  name: string;
  address: string;
  totalPrice?: number;
  seqCode?: number;
  status?: string;
  items?: InvoiceItem[];
}
