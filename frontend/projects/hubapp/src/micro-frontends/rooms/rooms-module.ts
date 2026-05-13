import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { EVENT_BUS_SERVICE} from '@ha/components-library/service';
import { EventBusService } from '@shared/service/event-bus/event-hub';

const routes = [
  {
    path: '',
    loadComponent: () =>
      import('./rooms/rooms.component').then((m) => m.RoomsComponent),
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
    },
  ]
})
export class RoomsModule {}
