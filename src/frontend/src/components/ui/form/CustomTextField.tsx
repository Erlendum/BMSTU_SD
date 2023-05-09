import React from 'react'
import { TextField } from '@material-ui/core'

type CustomTextFieldProps = {
	label: string
	name: string
	error: boolean
	helperText: string
	hide: boolean
	changeHandler: (event: React.ChangeEvent<HTMLInputElement>) => void
}

const CustomTextField = (props: CustomTextFieldProps) => {
	return (
		<TextField
			label={props.label}
			name={props.name}
			error={props.error}
			helperText={props.helperText}
			onChange={props.changeHandler}
			variant={'outlined'} //enables special material-ui styling
			size={'small'}
			margin={'dense'}
			type={props.hide ? 'password' : 'text'}
		/>
	)
}

export default CustomTextField
