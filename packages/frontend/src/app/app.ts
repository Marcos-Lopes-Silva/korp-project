import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { Sidebar } from './features/sidebar/sidebar.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, Sidebar],
  templateUrl: './app.html',
})
export class App {

}
