import React, { useState, useEffect, useContext } from "react";
import { makeStyles } from "@material-ui/core/styles";
import GridList from "@material-ui/core/GridList";
import GridListTile from "@material-ui/core/GridListTile";
import GridListTileBar from "@material-ui/core/GridListTileBar";
import IconButton from "@material-ui/core/IconButton";
import ShoppingCartOutlinedIcon from "@material-ui/icons/ShoppingCartOutlined";
import { getArt } from "../services/art";
import { RoutingContext, pagesMapping } from "../context/RoutingContext";
import selectedTileContext from "../context/SelectedTileContext";
import GalleryAppBar from "../components/GalleryAppBar";
import Paper from "@material-ui/core/Paper";

const useStyles = makeStyles((theme) => ({
  root: {
    display: "flex",
    flexWrap: "wrap",
    justifyContent: "space-around",
    overflow: "hidden",
    backgroundColor: theme.palette.background.paper,
  },
  appBar: {
    position: "relative",
  },
  layout: {
    width: "auto",
    marginLeft: theme.spacing(2),
    marginRight: theme.spacing(2),
    [theme.breakpoints.up(600 + theme.spacing(2) * 2)]: {
      width: 600,
      marginLeft: "auto",
      marginRight: "auto",
    },
  },
  paper: {
    marginTop: theme.spacing(3),
    marginBottom: theme.spacing(3),
    padding: theme.spacing(2),
    [theme.breakpoints.up(600 + theme.spacing(3) * 2)]: {
      marginTop: theme.spacing(6),
      marginBottom: theme.spacing(6),
      padding: theme.spacing(3),
    },
  } /*
  gridList: {
    width: 700,
    height: 800,
  },*/,
  icon: {
    color: "rgba(255, 255, 255, 0.54)",
  },
}));

const Gallery = () => {
  const { setPage } = useContext(RoutingContext);
  //const [selectedTile, setSelectedTile] = useState(null);

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

  const { setSelectedTile } = useContext(selectedTileContext);

  const handleClickBuy = (tile) => {
    setSelectedTile(tile);
    setPage(pagesMapping.buy);
  };

  return (
    <React.Fragment>
      <GalleryAppBar></GalleryAppBar>
      <main className={classes.layout}>
        <Paper className={classes.paper}>
          <GridList /*cellHeight={350} */ className={classes.gridList}>
            {art.map((tile) => (
              <GridListTile key={tile.small_file_url}>
                <img src={tile.small_file_url} alt={tile.title} />
                <GridListTileBar
                  title={tile.title}
                  subtitle={<span>by: {tile.artist}</span>}
                  actionIcon={
                    <IconButton
                      aria-label={`buy art ${tile.title}`}
                      className={classes.icon}
                      onClick={() => handleClickBuy(tile)}
                    >
                      <ShoppingCartOutlinedIcon />
                    </IconButton>
                  }
                />
              </GridListTile>
            ))}
          </GridList>
        </Paper>
      </main>
    </React.Fragment>
  );
};

export default Gallery;
