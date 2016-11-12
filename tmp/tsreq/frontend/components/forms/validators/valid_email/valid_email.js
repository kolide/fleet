"use strict";
var EMAIL_REGEX = /\S+@\S+\.\S+/;
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = function (email) {
    if (EMAIL_REGEX.test(email)) {
        return true;
    }
    return false;
};
