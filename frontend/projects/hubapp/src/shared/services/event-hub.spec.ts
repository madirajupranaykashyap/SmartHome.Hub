import { TestBed } from '@angular/core/testing';

import { EventHub } from '../../../../shared-library/src/lib/service/event-bus/event-hub';

describe('EventHub', () => {
  let service: EventHub;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(EventHub);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
