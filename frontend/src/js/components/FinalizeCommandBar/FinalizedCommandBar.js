import React from 'react';
import { connect } from 'react-redux';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';
import RippleButton from '../RippleButton';
import { CONFIG } from '../../utils/Utils';
import CopyToClipboard from 'react-copy-to-clipboard';

const buttonStyle = {
  'width': '95px',
  'height': '30px',
  'display': 'flex',
  'alignItems': 'center',
  'justifyContent': 'center'
}

class FinalizedCommandBar extends React.Component {

  constructor() {
    super();
    this.state = {
      show: true,
      cpyval: 'Copy Link',
      timoutFunc: null
    }
    this.getLink = this.getLink.bind(this);
    this.showLink = this.showLink.bind(this);
    this.getButton = this.getButton.bind(this);
    this.feedBackButton = this.feedBackButton.bind(this);
  }
  
  getLink() {
    if (this.state.show) {
      return (<input class="linkField" type="text" readOnly value={CONFIG.baseLocalURL + "/open/" + this.props.token}></input>)
    }
  }
  
  showLink() {
      if (!this.state.show) {
        this.setState({
          show: !this.state.show
        });  
      } else {

      }
  }

  feedBackButton() {
    const successStyle = {
      'fontSize': '14px'
    }

    this.setState({cpyval: <i class="material-icons" style={successStyle}>done</i>});

    this.state.timeoutFunc = setTimeout(function() {
        this.setState({cpyval: 'Copy Link'});
    }.bind(this), 2000);
  }

  componentWillUnmount()  {
    clearTimeout(this.state.timeoutFunc);
  }

  getButton() {
    if (this.state.show) {
      return (
        
        <CopyToClipboard text={CONFIG.baseLocalURL + "/#/open/" + this.props.token} onCopy={this.feedBackButton}>
          <button class="stylized green" style={buttonStyle}>{this.state.cpyval}</button>
        </CopyToClipboard>
        
      )
    } else {
      return <RippleButton cssClass="stylized green" customStyle={buttonStyle} action={this.showLink} label="Show Link"></RippleButton>
    }
  }

  render() {
    return (
      <div class="finalizedCommandBar">
        <ReactCSSTransitionGroup transitionName="field" transitionEnterTimeout={500} transitionLeaveTimeout={300}>
          {this.getLink()}
        </ReactCSSTransitionGroup>
        {this.getButton()}          
      </div>
    )
  } 
}

const mapStateToProps = (state) => {
  return {
    token: state.essentials.key
  }
} 

const mapDispatchToProps = (dispatch) => {
  return {
  
  }
}

const Finalized = connect(
  mapStateToProps,
  mapDispatchToProps
)(FinalizedCommandBar);

export default Finalized;