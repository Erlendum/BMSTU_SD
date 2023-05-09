import React from 'react'
import { TextField } from '@material-ui/core'

type CustomNumberInputProps = {
	label: string
	name: string
	error: boolean
	changeHandler: (event: React.ChangeEvent<HTMLInputElement>) => void
}

const CustomNumberInput = (props: CustomNumberInputProps) => {
	return (
		<TextField
			type={'number'}
			inputProps={{ inputMode: 'numeric', pattern: '[0-9]*' }}
			label={props.label}
			name={props.name}
			error={props.error}
			onChange={props.changeHandler}
			variant={'outlined'} //enables special material-ui styling
			size={'small'}
			margin={'dense'}
		/>
	)
}

export default CustomNumberInput
