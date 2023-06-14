import {
  STORE_ORDER_SUCCESS,
  STORE_ORDER_FAILURE,
  GET_ORDER_SUCCESS,
  GET_ORDER_FAILURE,
  GET_ALL_ORDERS_SUCCESS,
  GET_ALL_ORDERS_FAILURE,
  UPDATE_ORDER_STATUS_SUCCESS,
  UPDATE_ORDER_STATUS_FAILURE,
} from "./Action";

const initialState = {
  order: null,
  allOrders: [],
  error: null,
};

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case STORE_ORDER_SUCCESS:
      return {
        ...state,
        order: action.payload,
        error: null,
      };
    case STORE_ORDER_FAILURE:
    case GET_ORDER_FAILURE:
    case GET_ALL_ORDERS_FAILURE:
    case UPDATE_ORDER_STATUS_FAILURE:
      return {
        ...state,
        error: action.error,
      };
    case GET_ORDER_SUCCESS:
      return {
        ...state,
        order: action.payload,
        error: null,
      };
    case GET_ALL_ORDERS_SUCCESS:
      return {
        ...state,
        allOrders: action.payload,
        error: null,
      };
    case UPDATE_ORDER_STATUS_SUCCESS:
      return {
        ...state,
        order: action.payload,
        error: null,
      };
    default:
      return state;
  }
};

export default reducer;
