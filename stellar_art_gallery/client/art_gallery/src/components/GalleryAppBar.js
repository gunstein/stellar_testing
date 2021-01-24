import AppBar from "@material-ui/core/AppBar";
import { makeStyles } from "@material-ui/core/styles";
import Typography from "@material-ui/core/Typography";
import Toolbar from "@material-ui/core/Toolbar";

const useStyles = makeStyles((theme) => ({
  appBar: {
    position: "relative",
  },
}));

const GalleryAppBar = () => {
  const classes = useStyles();

  return (
    <AppBar position="absolute" color="default" className={classes.appBar}>
      <Toolbar>
        <Typography variant="h6" color="inherit" noWrap>
          StellarGallery
        </Typography>
      </Toolbar>
    </AppBar>
  );
};

export default GalleryAppBar;
