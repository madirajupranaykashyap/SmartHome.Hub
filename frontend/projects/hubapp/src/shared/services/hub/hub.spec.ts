import { TestBed } from '@angular/core/testing';

import { Hub } from './hub';

describe('Hub', () => {
  let service: Hub;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(Hub);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
