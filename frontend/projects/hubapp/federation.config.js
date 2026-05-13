const { withNativeFederation, shareAll} = require('@angular-architects/native-federation/config');

module.exports = withNativeFederation({
  name: 'hubapp',

  exposes: {
    './RoomsModule': './projects/hubapp/src/micro-frontends/rooms/rooms-module.ts',
    './DashboardModule': './projects/hubapp/src/micro-frontends/dashboard/dashboard-module.ts',
  },

  shared: {
    ...shareAll({ singleton: true, strictVersion: true, requiredVersion: 'auto' }),
    "@microsoft/signalr": {singleton: true, strictVersion: false, requiredVersion: 'auto'},
    "@angular/cdk": {singleton: true, strictVersion: false, requiredVersion: 'auto'},
  },

  skip: [
    'rxjs/ajax',
    'rxjs/fetch',
    'rxjs/testing',
    'rxjs/webSocket',
    // Add further packages you don't need at runtime
  ],

  // Please read our FAQ about sharing libs:
  // https://shorturl.at/jmzH0

  features: {
    // New feature for more performance and avoiding
    // issues with node libs. Comment this out to
    // get the traditional behavior:
    ignoreUnusedDeps: true,
  },
});
