import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';

interface JwtPayload {
  exp?: number;
}

const getJwtPayload = (token: string): JwtPayload | null => {
  const [, payload] = token.split('.');
  if (!payload) {
    return null;
  }

  try {
    const normalizedPayload = payload
      .replace(/-/g, '+')
      .replace(/_/g, '/')
      .padEnd(Math.ceil(payload.length / 4) * 4, '=');
    const decodedPayload = atob(normalizedPayload);

    return JSON.parse(decodedPayload) as JwtPayload;
  } catch {
    return null;
  }
};

const isJwtValid = (token: string): boolean => {
  const payload = getJwtPayload(token);
  if (!payload?.exp) {
    return false;
  }

  return payload.exp * 1000 > Date.now();
};

export const authGuard: CanActivateFn = (route, state) => {
  const router = inject(Router);

  const token = localStorage.getItem('token');

  if (token && token !== 'undefined' && token !== 'null' && isJwtValid(token)) {
    return true;
  }

  localStorage.removeItem('token');

  // Save requested URL
  localStorage.setItem(
    'redirect_url',
    state.url
  );

  return router.parseUrl('/auth/login');
};
