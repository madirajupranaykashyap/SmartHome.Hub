import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SharedLibrary } from './shared-library';

describe('SharedLibrary', () => {
  let component: SharedLibrary;
  let fixture: ComponentFixture<SharedLibrary>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SharedLibrary]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SharedLibrary);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
