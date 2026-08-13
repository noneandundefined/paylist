interface GUISwitchProps extends React.InputHTMLAttributes<HTMLInputElement> {}

const GUISwitch: React.FC<GUISwitchProps> = ({ className, ...props }) => {
	return (
		<div className="flex gap-2 items-center">
			<input type="checkbox" role="switch" className={`check-custom ${className || ''}`} {...props} />
		</div>
	);
};

export default GUISwitch;
