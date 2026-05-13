import { loadRemoteModule } from '@angular-architects/native-federation';
import { Routes } from '@angular/router';
import { authGuard } from './auth.guard';

export const routes: Routes = [
    {
        path: '',
        redirectTo: 'rooms',
        pathMatch: 'full',
    },
    {
        path: 'rooms',
        canActivate: [authGuard],
        loadChildren: () => 
            loadRemoteModule('hubapp', './RoomsModule')
            .then((m) => m.RoomsModule),
    },
    {
        path: 'auth',
        loadChildren: () =>
            import('../micro-frontends/auth/auth-module').then((m) => m.AuthModule),
            
    }
];
