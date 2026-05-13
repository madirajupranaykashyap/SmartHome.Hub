import { HttpInterceptorFn } from '@angular/common/http';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const token = localStorage.getItem('token');

  console.log('[+] Intercepting request to:', req.url);
  if (token) {

    req = req.clone({
      setHeaders: {
        Authorization: `Bearer ${token}`
      }
    });

  }
  return next(req);
};
