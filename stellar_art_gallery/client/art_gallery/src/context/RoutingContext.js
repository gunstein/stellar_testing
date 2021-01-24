import React, { useState, useMemo } from 'react'

export const pagesMapping = { 
  gallery: "gallery",
  buy: "buy"
};

export const RoutingContext = React.createContext({ page: pagesMapping.gallery });

export default function Router({ children }) {
  let urlPath = window.location.pathname.slice(1).toLowerCase()
  const [page, setPage] = useState(urlPath || pagesMapping.gallery);

  const value = useMemo(
    () => ({ page, setPage }), 
    [page, setPage]
  );
  
  return (
    <RoutingContext.Provider value={value}>
      {children}
    </RoutingContext.Provider>
  )
}