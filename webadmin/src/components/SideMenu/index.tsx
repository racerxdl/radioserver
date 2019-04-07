import Drawer from "@material-ui/core/Drawer";
import Hidden from "@material-ui/core/Hidden";
import {Component, default as React} from "react";
import {RouteComponentProps, withRouter} from "react-router";
import MenuButton from '../MenuButton';

type SideMenuProps = {
  classes: any
  theme: any,
  mobileOpen: boolean,
  handleDrawerToggle?: () => void,
}

type AllMenuProps = RouteComponentProps & SideMenuProps;

class SideMenu extends Component<AllMenuProps> {
  render() {
    const menuContainer = (
      <div>
        <div className={this.props.classes.toolbar}/>
        <MenuButton icon={'home'} content={'Home'} onClick={() => this.props.history.push('/')}/>
        <br/>
        <MenuButton icon={'dashboard'} content={'Dashboard'} onClick={() => this.props.history.push('/dashboard')}/>
        <br/>
        <MenuButton icon={'settings'} content={'Settings'} onClick={() => this.props.history.push('/settings')}/>
      </div>
    );

    return (
      <nav className={this.props.classes.drawer}>
        <Hidden smUp implementation="css">
          <Drawer
            variant="temporary"
            anchor={this.props.theme.direction === 'rtl' ? 'right' : 'left'}
            open={this.props.mobileOpen}
            onClose={this.props.handleDrawerToggle}
            classes={{
              paper: this.props.classes.drawerPaper,
            }}>
            {menuContainer}
          </Drawer>
        </Hidden>
        <Hidden xsDown implementation="css">
          <Drawer
            classes={{
              paper: this.props.classes.drawerPaper,
            }}
            variant="permanent"
            open>
            {menuContainer}
          </Drawer>
        </Hidden>
      </nav>
    );
  };
}

export default withRouter(SideMenu);
