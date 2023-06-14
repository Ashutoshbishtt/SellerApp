import React from "react";
import { Provider } from "react-redux";
import { createStore } from "redux";
import store from "./Components/Redux/Store";

import rootReducer from "./Components/Redux/Reducers";
import ExcelDragAndDrop from "./Components/Pages/DragAndDrop";

// Create your Redux store

function App() {
  return (
    <Provider store={store}>
      <div className="App">
        <ExcelDragAndDrop />
      </div>
    </Provider>
  );
}

export default App;
