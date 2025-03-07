import * as React from "react";
import { createRoot } from "react-dom/client";
import { App } from "./Components/App/App";

const root = createRoot(document.getElementById("app")!);

root.render(<App />);
