import React from 'react';

export default class RippleButton extends React.Component {

	constructor() {
		super();
		this.action = this.action.bind(this);
	}

	action(e) {

		this.ink.classList.remove("animate");

		if (!this.ink.clientHeight && !this.ink.clientWidth) {
			const d = Math.max(this.button.offsetWidth, this.button.offsetHeight);
			this.ink.style.height = d + "px";
			this.ink.style.width = d + "px";
		}

		const x = e.pageX - this.button.offsetLeft - (this.ink.clientWidth / 2);
		const y = e.pageY - this.button.offsetTop - (this.ink.clientHeight / 2);
	
		const inkColor = (this.props.rippleColor) ? this.props.rippleColor : 'rgba(255, 255, 255, 0.4)';

		this.ink.style.top = y + "px";
		this.ink.style.left = x + "px";
		this.ink.style.background = inkColor;
		this.ink.classList.add("animate");

		this.props.action(e, this.props.actionAttr);	
	}

	render() {
		return (
			<button ref={(c) => this.button = c} class={this.props.cssClass} style={this.props.customStyle} onClick={(e) => this.action(e)}>
				<span class="ink" ref={(c) => this.ink = c}></span>
				{this.props.label}
			</button>
		)
	}

}
