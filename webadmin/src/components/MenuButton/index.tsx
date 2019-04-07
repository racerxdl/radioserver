import Button from "@material-ui/core/Button";
import * as React from "react";
import withStyles from "@material-ui/core/styles/withStyles";
import Icon from '@material-ui/core/Icon';

type MenuButtonProps = {
  icon: any,
  content: any,
  classes: any,
  onClick?: () => void,
}

const MenuButton = (props: MenuButtonProps) => {
  return (
    <Button
      focusRipple={true}
      classes={{root: props.classes.button, label: props.classes.label}}
      onClick={props.onClick}>
      <Icon
        className={props.classes.icon}>
        {props.icon}
      </Icon>
      {props.content}
    </Button>
  );
};

export default withStyles((theme: any) => ({
  button: {
    width: '100%',
    height: '96px'
  },
  label: {
    flexDirection: 'column'
  },
  icon: {
    fontSize: '32px !important',
    marginBottom: theme.spacing.unit
  }
}))
(MenuButton);
