import React, { FC, useEffect, useState } from 'react'
import InstrumentItem from '../instrument-item/InstrumentItem'
import { IInstrument } from '../../../types/instrument.interface'
import { UserService } from '../../../services/user.service'
import { toast, ToastContainer } from 'react-toastify'
import { InstrumentService } from '../../../services/instrument.service'
import { ComparisonListService } from '../../../services/comparisonList.service'
import styles from '../instrument-item-table/InstrumentItemTable.module.scss'

const InstrumentItemTable: FC<{
	instrument: IInstrument
}> = ({ instrument }) => {
	const [error, setError] = useState('no error')

	useEffect(() => {
		if (error != 'no error') {
			toast.error('ERROR ' + error, {
				position: toast.POSITION.BOTTOM_LEFT
			})
		}
	}, [error])

	const handleDeleteComparisonList = async () => {
		setError('no error')
		let isError = false
		await ComparisonListService.deleteInstrument(instrument.InstrumentId).catch(
			error => {
				isError = true
				if (error.response) {
					setError(error.response.data.Error)
				}
			}
		)
		if (!isError) {
			toast.success('Instrument was successfully deleted', {
				position: toast.POSITION.BOTTOM_LEFT
			})
		}
		console.log(error)
	}

	return (
		<tr>
			<td className={styles.textLeft}>
				<img src={instrument.Img} alt={instrument.Name} />
			</td>
			<td className={styles.textLeft}>{instrument.Name}</td>
			<td className={styles.textLeft}>{instrument.Brand}</td>
			<td className={styles.textLeft}>{instrument.Material}</td>
			<td className={styles.textLeft}>{instrument.Type}</td>
			<td className={styles.textLeft}>
				{new Intl.NumberFormat('ru-RU', {
					style: 'currency',
					currency: 'RUB',
					maximumFractionDigits: 0
				}).format(instrument.Price)}
			</td>
			<td className={styles.textLeft}>
				<a href='javascript:void(0);' onClick={handleDeleteComparisonList}>
					<i className='fa fa-trash'></i>
				</a>
			</td>
		</tr>
	)
}

export default InstrumentItemTable
