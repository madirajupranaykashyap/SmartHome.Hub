import { CommonModule } from '@angular/common';
import { Component, DestroyRef, inject, OnInit } from '@angular/core';
import { HubLoginModule } from '@ha/components-library/pages/hub-login';
import { EVENT_BUS_SERVICE } from '@ha/components-library/service';
import { HubService, LoginResponse } from '../../../../../hubapp/src/shared/services/hub/hub';
import { Router } from '@angular/router';

const getAccessToken = (response: LoginResponse | string): string | null => {
  if (typeof response === 'string') {
    return response;
  }

  return response.access_token
    ?? response.accessToken
    ?? response.token
    ?? response.jwt
    ?? null;
};

@Component({
  selector: 'app-login',
  imports: [CommonModule, HubLoginModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.less',
})
export class LoginComponent implements OnInit {
  private eventBus = inject(EVENT_BUS_SERVICE);
  private hubService = inject(HubService);
  private destroyRef = inject(DestroyRef);
  private router = inject(Router);

  ngOnInit() {
    this.eventBus.on('hub-login-form', this.destroyRef).subscribe((data) => {
      console.log('[+] Received login data:', data);
      this.hubService.login(data.username, data.password).subscribe({
        next: (response) => {
          const token = getAccessToken(response);

          if (!token) {
            console.error('Login response did not include a JWT token:', response);
            return;
          }

          localStorage.setItem('token', token);
          console.log('Login successful:', response);
          const redirectUrl = localStorage.getItem('redirect_url');

          if (redirectUrl) {
            localStorage.removeItem('redirect_url');

            this.router.navigateByUrl(redirectUrl);
          } else {
            // fallback page
            this.router.navigate(['/']);
          }
        },
      });
    });
  }
}
