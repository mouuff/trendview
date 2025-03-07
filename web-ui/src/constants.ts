declare const __DEV__: boolean;

export const RELATED_RATING_CUTOFF = 40;
export const BUY_THRESHOLD = 60;
export const DONT_BUY_THRESHOLD = 40;

export const SERVER_BASE_URL = __DEV__ ? "http://localhost:8081/" : "/";
