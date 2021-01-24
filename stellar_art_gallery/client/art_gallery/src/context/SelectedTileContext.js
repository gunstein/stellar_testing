import { createContext } from "react";

const selectedTileContext = createContext({
  selectedTile: null,
  setSelectedTile: (auth) => {}
});

export default selectedTileContext;