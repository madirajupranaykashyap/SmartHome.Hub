import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { HubUserRoom } from '@ha/components-library/models';
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
  _rooms$ = new BehaviorSubject<DataState<HubUserRoom[]>>({ state: 'waitingForData' });
  readonly rooms$ = this._rooms$.asObservable();

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

  getAllRooms(): Observable<HubUserRoom[]> {
    const url = this.baseUrl + `/api/rooms`;
    this._rooms$.next({ state: 'waitingForData' });
    return this.http.get<HubUserRoom[]>(url).pipe(
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
}
