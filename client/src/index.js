import React from "react";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import store from "./Components/Redux/Store";
import App from "./App";

const rootElement = document.getElementById("root");

createRoot(rootElement).render(
  <Provider store={store}>
    <App />
  </Provider>
);
