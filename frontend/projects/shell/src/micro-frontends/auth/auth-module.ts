import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes } from '@angular/router';
import { EVENT_BUS_SERVICE } from '@ha/components-library/service';
import { EventBusService } from '@shared/service/event-bus/event-hub';

const routes: Routes = [
  {
    path: '',
    redirectTo: 'login',
    pathMatch: 'full',
  },
  {
    path: 'login',
    loadComponent: () =>
      import('./login/login.component').then((m) => m.LoginComponent),
  },
  {
    path: 'register',
    loadComponent: () =>
      import('./register/register.component').then((m) => m.RegisterComponent),
  }
]
@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    RouterModule.forChild(routes)
  ],
  providers: [
    {
      provide: EVENT_BUS_SERVICE,
      useClass: EventBusService
    }
  ]
})
export class AuthModule {}
