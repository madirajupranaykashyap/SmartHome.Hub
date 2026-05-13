import { DestroyRef, Injectable } from '@angular/core';
import { EventMap, IEventBusService } from '@ha/components-library/service';
import { Observable, ReplaySubject } from 'rxjs';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Injectable()
export class EventBusService implements IEventBusService { 
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    private channels = new Map<keyof EventMap, ReplaySubject<any>>();

    private id = Math.random();

    constructor() {
        //   console.log('EventBus instance ID:', this.id);
    }

    // 🔥 Centralized typed getter
    private getSubject<K extends keyof EventMap>(
        event: K
    ): ReplaySubject<EventMap[K]> {

        if (!this.channels.has(event)) {
            this.channels.set(event, new ReplaySubject<EventMap[K]>(1));
        }

        return this.channels.get(event)! as ReplaySubject<EventMap[K]>;
    }

    emit<K extends keyof EventMap>(
        event: K,
        payload: EventMap[K]
    ): void {
        this.getSubject(event).next(payload);
    }
    
    
    on<K extends keyof EventMap>(
        event: K,
        destroyRef: DestroyRef
    ): Observable<EventMap[K]> {

        return this.getSubject(event)
            .asObservable()
            .pipe(takeUntilDestroyed(destroyRef));
    }
    
    clear<K extends keyof EventMap>(event: K): void {
        this.channels.get(event)?.complete();
        this.channels.delete(event);
    }
}
