import React, { useState, useContext } from "react";
import Gallery from "./pages/Gallery";
import Buy from "./pages/Buy";
import { pagesMapping, RoutingContext } from "./context/RoutingContext"
import selectedTileContext from "./context/SelectedTileContext";
import './App.css';

export default function App() {
  const { page } = useContext(RoutingContext)
  const [selectedTile, setSelectedTile] = useState(null);
  return (
    <selectedTileContext.Provider value={{ selectedTile, setSelectedTile }}>
      <React.Fragment>
        {(pagesMapping.gallery === page) && <Gallery />}
        {(pagesMapping.buy === page) && <Buy />}
      </React.Fragment>
    </selectedTileContext.Provider>
  );
}
