import React, { FC, useEffect, useState } from 'react'
import Layout from '../../ui/layout/Layout'
import { useQuery } from 'react-query'
import { UserService } from '../../../services/user.service'
import InstrumentItem from '../../ui/instrument-item/InstrumentItem'
import styles from './ComparisonList.module.scss'
import { ToastContainer, toast } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'
import { ComparisonListService } from '../../../services/comparisonList.service'
import InstrumentItemTable from '../../ui/instrument-item-table/InstrumentItemTable'

const ComparisonList: FC = () => {
	const [error, setError] = useState('no error')
	let instrumentId = 0
	const [updateQuery, setUpdateQuery] =
		useState(false)
	const { data: comparisonListInstruments, isLoading } = useQuery(
		['comparisonListInstruments', updateQuery],
		() => UserService.getComparisonList(),
		{
			select: ({ comparisonListInstruments }) => comparisonListInstruments
		}
	)

	const { data: comparisonList } = useQuery(
		['comparisonList', updateQuery],
		() => UserService.getComparisonList(),
		{
			select: ({ comparisonList }) => comparisonList
		}
	)

	const displayComparisonListInstruments = comparisonListInstruments?.map(
		instrument => {
			return (
				<InstrumentItemTable
					updateQuery={updateQuery}
					setUpdateQuery={setUpdateQuery}
					instrument={instrument}
					key={instrument.InstrumentId}
				/>
			)
		}
	)

	return (
		<Layout title='Comparison List'>
			<div className={styles.text}>Amount: {comparisonList?.Amount}</div>
			<div className={styles.text}>
				Total Price: {comparisonList?.TotalPrice}
			</div>
			<table className={styles.table}>
				<thead>
					<tr>
						<th className={styles.textLeft}>Img</th>
						<th className={styles.textLeft}>Name</th>
						<th className={styles.textLeft}>Brand</th>
						<th className={styles.textLeft}>Material</th>
						<th className={styles.textLeft}>Type</th>
						<th className={styles.textLeft}>Price</th>
						<th className={styles.textLeft}>Links</th>
					</tr>
				</thead>
				<tbody className='table-hover'>
					{displayComparisonListInstruments}
				</tbody>
			</table>
			<ToastContainer />
		</Layout>
	)
}

export default ComparisonList
