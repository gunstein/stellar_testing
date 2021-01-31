import React, { useState, useEffect } from "react";
import { makeStyles } from "@material-ui/core/styles";
import GridList from "@material-ui/core/GridList";
import GridListTile from "@material-ui/core/GridListTile";
import GridListTileBar from "@material-ui/core/GridListTileBar";
import ListSubheader from "@material-ui/core/ListSubheader";
import IconButton from "@material-ui/core/IconButton";
import ShoppingCartOutlinedIcon from "@material-ui/icons/ShoppingCartOutlined";
import { getArt } from "../services/art";

const useStyles = makeStyles((theme) => ({
  root: {
    display: "flex",
    flexWrap: "wrap",
    justifyContent: "space-around",
    overflow: "hidden",
    backgroundColor: theme.palette.background.paper,
  },
  gridList: {
    width: 700,
    height: 800,
  },
  icon: {
    color: "rgba(255, 255, 255, 0.54)",
  },
}));

export default function GalleryTitlebarGridList({ handleOrder }) {
  const classes = useStyles();
  const [art, setArt] = useState([]);
  useEffect(() => {
    let mounted = true;
    getArt().then((items) => {
      if (mounted) {
        //console.log(items.data)
        setArt(items.data);
      }
    });
    return () => (mounted = false);
  }, []);

  const handleClickOpen = (tile) => {
    handleOrder(tile);
  };

  return (
    <div className={classes.root}>
      <GridList cellHeight={350} className={classes.gridList}>
        <GridListTile key="Subheader" cols={2} style={{ height: "auto" }}>
          <ListSubheader component="div" style={{ fontSize: "50px" }}>
            StellarGallery
          </ListSubheader>
        </GridListTile>
        {art.map((tile) => (
          <GridListTile key={tile.small_file_url}>
            <img
              src={`${process.env.PUBLIC_URL} + {tile.small_file_url}`}
              alt={tile.title}
            />
            <GridListTileBar
              title={tile.title}
              subtitle={<span>by: {tile.artist}</span>}
              actionIcon={
                <IconButton
                  aria-label={`buy art ${tile.title}`}
                  className={classes.icon}
                  onClick={() => handleClickOpen(tile)}
                >
                  <ShoppingCartOutlinedIcon />
                </IconButton>
              }
            />
          </GridListTile>
        ))}
      </GridList>
    </div>
  );
}
