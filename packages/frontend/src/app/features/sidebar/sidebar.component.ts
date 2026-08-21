import { Component, inject } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { ServiceStatus, StatusService } from '../../core/services/status.service';
import { CommonModule } from '@angular/common';
import { forkJoin, interval, map, Observable, startWith, switchMap } from 'rxjs';
import { RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-sidebar',
  imports: [CommonModule, RouterLink, RouterLinkActive],
  templateUrl: './sidebar.component.html',
})
export class Sidebar {

  private statusService = inject(StatusService)

  services = toSignal(
    interval(10000).pipe(
      startWith(0),
      switchMap(() =>
        forkJoin([
          this.statusService.checkHealth('Estoque', '/stock'),
          this.statusService.checkHealth('Faturamento', '/billing'),
          this.statusService.checkHealth('Gateway', ''),
        ]) as Observable<ServiceStatus[]>
      )
    ),
    { initialValue: [] as ServiceStatus[] }
  );
}

