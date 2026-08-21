import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, catchError, of, map } from 'rxjs';
import { environment } from '../../environments/environment';

export interface ServiceStatus {
  name: string;
  online: boolean;
}

@Injectable({ providedIn: 'root' })
export class StatusService {
  private gatewayUrl = environment.apiUrl;

  constructor(private http: HttpClient) { }

  checkHealth(serviceName: string, path: string): Observable<ServiceStatus> {
    return this.http.get(`${this.gatewayUrl}${path}/health`).pipe(
      map(() => ({ name: serviceName, online: true })),
      catchError(() => of({ name: serviceName, online: false }))
    );
  }
}
