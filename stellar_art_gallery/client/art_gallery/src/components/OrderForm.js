import React from "react";
import Button from "@material-ui/core/Button";
import TextField from "@material-ui/core/TextField";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogTitle from "@material-ui/core/DialogTitle";

export default function FormDialog({ handleCloseToParent, tile }) {
  //const [open, setOpen] = React.useState(false);

  const handleClose = () => {
    handleCloseToParent();
  };

  return (
    <Dialog
      open={true}
      onClose={handleClose}
      aria-labelledby="form-dialog-title"
    >
      <DialogTitle id="form-dialog-title">
        Buy {tile.title} by {tile.artist}
      </DialogTitle>
      <DialogContent>
        <DialogContentText>
          Please enter email and press buy. Then paymentinfo will show up. When
          payment is received you will receive downloadinfo per email.
        </DialogContentText>
        <TextField
          autoFocus
          margin="dense"
          id="name"
          label="Email Address"
          type="email"
          fullWidth
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} color="primary">
          Cancel
        </Button>
        <Button onClick={handleClose} color="primary">
          Buy
        </Button>
      </DialogActions>
    </Dialog>
  );
}
