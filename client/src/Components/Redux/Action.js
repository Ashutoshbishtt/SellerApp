export const STORE_ORDER_SUCCESS = "STORE_ORDER_SUCCESS";
export const STORE_ORDER_FAILURE = "STORE_ORDER_FAILURE";
export const GET_ORDER_SUCCESS = "GET_ORDER_SUCCESS";
export const GET_ORDER_FAILURE = "GET_ORDER_FAILURE";
export const GET_ALL_ORDERS_SUCCESS = "GET_ALL_ORDERS_SUCCESS";
export const GET_ALL_ORDERS_FAILURE = "GET_ALL_ORDERS_FAILURE";
export const UPDATE_ORDER_STATUS_SUCCESS = "UPDATE_ORDER_STATUS_SUCCESS";
export const UPDATE_ORDER_STATUS_FAILURE = "UPDATE_ORDER_STATUS_FAILURE";

export const storeOrder = orderPayload => async dispatch => {
  try {
    const response = await fetch("http://localhost:8080/storeOrder", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(orderPayload),
    });
    const data = await response.json();

    if (response.ok) {
      dispatch({
        type: STORE_ORDER_SUCCESS,
        payload: data,
      });
    } else {
      dispatch({
        type: STORE_ORDER_FAILURE,
        error: data,
      });
    }
  } catch (error) {
    dispatch({
      type: STORE_ORDER_FAILURE,
      error: error.message,
    });
  }
};

export const getOrder = orderId => async dispatch => {
  try {
    const response = await fetch(
      `http://localhost:8080/getOrder?id=${orderId}`
    );
    const data = await response.json();

    if (response.ok) {
      dispatch({
        type: GET_ORDER_SUCCESS,
        payload: data,
      });
    } else {
      dispatch({
        type: GET_ORDER_FAILURE,
        error: data,
      });
    }
  } catch (error) {
    dispatch({
      type: GET_ORDER_FAILURE,
      error: error.message,
    });
  }
};

export const getAllOrders = params => async dispatch => {
  try {
    const queryParams = new URLSearchParams(params).toString();
    const response = await fetch(
      `http://localhost:8080/getAllOrders?${queryParams}`
    );
    const data = await response.json();

    if (response.ok) {
      dispatch({
        type: GET_ALL_ORDERS_SUCCESS,
        payload: data,
      });
    } else {
      dispatch({
        type: GET_ALL_ORDERS_FAILURE,
        error: data,
      });
    }
  } catch (error) {
    dispatch({
      type: GET_ALL_ORDERS_FAILURE,
      error: error.message,
    });
  }
};

export const updateOrderStatus = orderStatusPayload => async dispatch => {
  try {
    const response = await fetch("http://localhost:8080/updateOrderStatus", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(orderStatusPayload),
    });
    const data = await response.json();

    if (response.ok) {
      dispatch({
        type: UPDATE_ORDER_STATUS_SUCCESS,
        payload: data,
      });
    } else {
      dispatch({
        type: UPDATE_ORDER_STATUS_FAILURE,
        error: data,
      });
    }
  } catch (error) {
    dispatch({
      type: UPDATE_ORDER_STATUS_FAILURE,
      error: error.message,
    });
  }
};
