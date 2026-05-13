import { Routes } from '@angular/router';

export const routes: Routes = [
    {
        path: '',
        redirectTo: 'rooms',
        pathMatch: 'full',
    },
    { 
        path: 'rooms', 
        loadChildren: () => 
            import('../micro-frontends/rooms/rooms-module').then(m => m.RoomsModule) 
    },
    { 
        path: 'dashboard', 
        loadChildren: () => 
            import('../micro-frontends/dashboard/dashboard-module').then(m => m.DashboardModule) 
    }
];
