import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Room, UserRoom } from '@ha/components-library/models';
import { DataState } from '@ha/components-library/shared';
import { BehaviorSubject, Observable, tap } from 'rxjs';

export interface LoginResponse {
  access_token?: string;
  accessToken?: string;
  token?: string;
  jwt?: string;
}

@Injectable({
  providedIn: 'root',
})
export class HubService {
  // baseUrl = "";
  baseUrl = 'http://localhost:8080';

  private http = inject(HttpClient);
  _rooms$ = new BehaviorSubject<DataState<UserRoom[]>>({ state: 'waitingForData' });
  readonly rooms$ = this._rooms$.asObservable();
  _allRooms$ = new BehaviorSubject<DataState<Room[]>>({ state: 'waitingForData' });
  readonly allRooms$ = this._allRooms$.asObservable();

  // login

  login(username: string, password: string): Observable<LoginResponse | string> {

    console.log('[+] Attempting login with username:', username);
    console.log('[+] Attempting login with password:', password);
    const body = new URLSearchParams();

    body.set('username', username);
    body.set('password', password);

    return this.http.post(
      `${this.baseUrl}/api/auth/login`,
      body.toString(),
      {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      }
    );
  }

  getAllRooms(): Observable<UserRoom[]> {
    const url = this.baseUrl + `/api/rooms`;
    this._rooms$.next({ state: 'waitingForData' });
    return this.http.get<UserRoom[]>(url).pipe(
      tap({
        next: (data) => {
          this._rooms$.next({
            state: 'hasData',
            data,
          });
        },
        error: (error) => {
          this._rooms$.next({
            state: 'hasError',
            error,
          });
        }
      })
    );
  }

  getRoomCatalog(): Observable<Room[]> {
    const url = this.baseUrl + `/api/rooms/catalog`;
    this._allRooms$.next({ state: 'waitingForData' });
    return this.http.get<Room[]>(url).pipe(
      tap({
        next: (data) => {
          this._allRooms$.next({
            state: 'hasData',
            data,
          });
        },
        error: (error) => {
          this._allRooms$.next({
            state: 'hasError',
            error,
          });
        }
      })
    );
  }

  createUserRoom(roomGuid: string): Observable<UserRoom> {
    const url = this.baseUrl + `/api/rooms`;
    return this.http.post<UserRoom>(url, { room_guid: roomGuid }).pipe(
      tap(() => {
        this.getAllRooms().subscribe();
      })
    );
  }

  deleteUserRoom(id: string): Observable<void> {
    const url = this.baseUrl + `/api/rooms/${id}`;
    return this.http.delete<void>(url).pipe(
      tap(() => {
        this.getAllRooms().subscribe();
      })
    );
  }
}
