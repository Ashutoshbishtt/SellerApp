import React, { useState } from "react";
import { read, utils } from "xlsx";
import { useDropzone } from "react-dropzone";
import { useDispatch } from "react-redux";
import { storeOrder } from "../Redux/Action";
import "./DragAndDrop.css";

function ExcelDragAndDrop() {
  const dispatch = useDispatch();
  const [isLoading, setIsLoading] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);
  const [isError, setIsError] = useState(false);

  const handleDrop = async acceptedFiles => {
    const file = acceptedFiles[0];
    const reader = new FileReader();

    setIsLoading(true);
    setIsSuccess(false);
    setIsError(false);

    reader.onload = async e => {
      const data = new Uint8Array(e.target.result);
      const workbook = read(data, { type: "array" });
      const sheetName = workbook.SheetNames[0];
      const worksheet = workbook.Sheets[sheetName];
      const jsonData = utils.sheet_to_json(worksheet, { header: 1 });

      const orders = [];
      for (let i = 1; i < jsonData.length; i++) {
        const row = jsonData[i];
        const orderData = {
          id: String(row[0]),
          status: row[1],
          items: [],
          total: 0,
          currencyUnit: row[7],
        };

        const item = {
          id: String(row[2]),
          description: row[3],
          price: parseFloat(row[4]),
          quantity: parseInt(row[5]),
        };

        const existingOrder = orders.find(o => o.id === orderData.id);
        if (existingOrder) {
          existingOrder.items.push(item);
          existingOrder.total += item.price * item.quantity;
        } else {
          orderData.items.push(item);
          orderData.total = item.price * item.quantity;
          orders.push(orderData);
        }
      }

      try {
        // Simulating an asynchronous save operation
        const result = await dispatch(storeOrder(orders));

        if (result.error) {
          setIsError(true);
          setIsSuccess(false);
          setIsLoading(false);
        } else {
          setIsError(false);
          setIsSuccess(true);
          setIsLoading(false);
        }
      } catch (error) {
        setIsError(true);
        setIsSuccess(false);
        setIsLoading(false);
      }

      console.log(orders);
    };

    reader.onerror = () => {
      setIsLoading(false);
      setIsError(true);
    };

    reader.readAsArrayBuffer(file);
  };

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop: handleDrop,
    multiple: false,
    accept: ".xlsx",
  });

  return (
    <div
      {...getRootProps()}
      className={`drag-area ${isDragActive ? "drag-active" : ""}`}
    >
      <input {...getInputProps()} />
      {isLoading ? (
        <div className="loading-spinner" />
      ) : (
        <>
          {isSuccess ? (
            <p className="success-message">
              File saved in the database successfully!
            </p>
          ) : (
            <>
              {isError ? (
                <p className="error-message">
                  Error occurred. Please select another file.
                </p>
              ) : (
                <p className="drop-text">
                  {isDragActive
                    ? "Drop the Excel file here..."
                    : "Drag and drop an Excel file here, or click to select a file"}
                </p>
              )}
            </>
          )}
        </>
      )}
    </div>
  );
}

export default ExcelDragAndDrop;
