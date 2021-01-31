import React from "react";
import Dialog from "@material-ui/core/Dialog";
import CloseIcon from "@material-ui/icons/Close";
import Fab from "@material-ui/core/Fab";
import GalleryAppBar from "../components/GalleryAppBar";
import Slide from "@material-ui/core/Slide";
import Paper from "@material-ui/core/Paper";
import { makeStyles } from "@material-ui/core/styles";

export default function ImageDialog({
  handleCloseToParent,
  bigImageUrl,
  alternativeBigImageUrl,
}) {
  const useStyles = makeStyles((theme) => ({
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
    },
  }));

  const classes = useStyles();

  const handleClose = () => {
    handleCloseToParent();
  };

  const addAlternativeImgSrc = (ev) => {
    ev.target.src = alternativeBigImageUrl;
  };

  const Transition = React.forwardRef(function Transition(props, ref) {
    return <Slide direction="up" ref={ref} {...props} />;
  });

  return (
    <Dialog
      fullScreen
      open={true}
      onClose={handleClose}
      TransitionComponent={Transition}
    >
      <main className={classes.layout}>
        <GalleryAppBar></GalleryAppBar>
        <Paper className={classes.paper}>
          <Fab color="primary" aria-label="Close" onClick={() => handleClose()}>
            <CloseIcon />
          </Fab>
          <img
            onError={addAlternativeImgSrc}
            src={bigImageUrl}
            alt={bigImageUrl}
            style={{ width: "100%" }}
          />
        </Paper>
      </main>
    </Dialog>
  );
}
