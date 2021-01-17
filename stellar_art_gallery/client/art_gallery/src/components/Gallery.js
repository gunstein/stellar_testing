import React, { useState } from "react";
import OrderForm from "./OrderForm.js";
import GalleryTitlebarGridList from "./GalleryTitlebarGridList";

const Gallery = () => {
  const [orderMode, setOrderMode] = useState(0);
  const [selectedTile, setSelectedTile] = useState(null);

  const handleOrder = (tile) => {
    //console.log(tile);
    setSelectedTile(tile);
    setOrderMode(1);
  };

  const handleCloseOrderDialog = () => {
    setOrderMode(0);
  };

  return (
    <React.Fragment>
      <GalleryTitlebarGridList
        handleOrder={handleOrder}
      ></GalleryTitlebarGridList>
      {orderMode ? (
        <OrderForm
          handleCloseToParent={handleCloseOrderDialog}
          tile={selectedTile}
        ></OrderForm>
      ) : null}
    </React.Fragment>
  );
};

export default Gallery;
