/*
 * Copyright 2016-present, Kolide, Inc.
 * All rights reserved.
 *
 * @flow
 */

import { Reactor } from 'nuclear-js';

/**
 * Dispatcher is the NuclearJS implementation of the Flux pattern dispatcher.
 *
 * @exports Dispatcher
 */
var Dispatcher = new Reactor({
  debug: true
});

export default Dispatcher;