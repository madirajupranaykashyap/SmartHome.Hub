import { Component, DestroyRef, inject, OnInit } from '@angular/core';
import { HubService } from '../../../shared/services/hub/hub';
import { Observable } from 'rxjs';
import { DataState } from '@ha/components-library/shared';
import { Room, UserRoom } from '@ha/components-library/models';
import { CommonModule } from '@angular/common';
import { RoomsModule } from '@ha/components-library/pages/rooms';
import { EVENT_BUS_SERVICE } from '@ha/components-library/service';

@Component({
  selector: 'app-rooms',
  imports: [
    CommonModule,
    RoomsModule
],
  templateUrl: './rooms.component.html',
  styleUrl: './rooms.component.less',
})
export class RoomsComponent implements OnInit {
  private hubService = inject(HubService);
  private eventBus = inject(EVENT_BUS_SERVICE);
  private destroyRef = inject(DestroyRef);
  
  userRooms$:  Observable<DataState<UserRoom[]>> = this.hubService.rooms$;
  allRooms$:  Observable<DataState<Room[]>> = this.hubService.allRooms$;

  ngOnInit(): void {
    this.hubService.getAllRooms().subscribe();
    this.hubService.getRoomCatalog().subscribe();
    this.eventBus.on('room-create', this.destroyRef).subscribe((data) => {
      this.hubService.createUserRoom(data.id).subscribe();
    });
    this.eventBus.on('room-delete', this.destroyRef).subscribe((data) => {
      this.hubService.deleteUserRoom(data.id).subscribe();
    });
    console.log('[+] RoomsComponent initialized, fetching rooms...');
    console.log('[+] Current rooms state:', this.hubService.rooms$);
  }
}
