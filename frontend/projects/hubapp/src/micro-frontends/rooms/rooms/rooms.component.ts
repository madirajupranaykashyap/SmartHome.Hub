import { Component, inject, OnInit } from '@angular/core';
import { HubRoomModule } from '@ha/components-library/pages/hub-room';
import { HubService } from '../../../shared/services/hub/hub';
import { Observable } from 'rxjs';
import { DataState } from '@ha/components-library/shared';
import { HubUserRoom } from '@ha/components-library/models';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-rooms',
  imports: [
    CommonModule,
    HubRoomModule,
],
  templateUrl: './rooms.component.html',
  styleUrl: './rooms.component.less',
})
export class RoomsComponent implements OnInit {
  private hubService = inject(HubService);
  
  userRooms$:  Observable<DataState<HubUserRoom[]>> = this.hubService.rooms$;

  ngOnInit(): void {
    this.hubService.getAllRooms().subscribe();
  }
}
