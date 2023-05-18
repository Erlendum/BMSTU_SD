import React, { FC, useEffect, useState } from 'react'
import { useQuery } from 'react-query'
import { UserService } from '../../../services/user.service'
import { toast, ToastContainer } from 'react-toastify'
import { OrderService } from '../../../services/order.service'
import InstrumentItemTable from '../../ui/instrument-item-table/InstrumentItemTable'
import Layout from '../../ui/layout/Layout'
import styles from '../orders/Orders.module.scss'
import OrderItem from '../../ui/order-item/OrderItem'

const Orders: FC = () => {
	const [error, setError] = useState('no error')
	let isAdmin =
		UserService.getCurrentIsAdmin() != null &&
		UserService.getCurrentIsAdmin() == 'true'
	const [updateQuery, setUpdateQuery] =
		useState(false)
	const { data: orders, isLoading } = useQuery(
		['orders', updateQuery],
		() => !isAdmin ? OrderService.getList() : OrderService.getListForAll(),
		{
			select: ({ orders }) => orders
		}
	)

	useEffect(() => {
		if (error !== 'no error') {
			toast.error('ERROR ' + error, {
				position: toast.POSITION.BOTTOM_LEFT
			})
		}
	}, [error])
	const handleCheckout = async (): Promise<void> => {
		let id: string = ''
		setError('no error')
		let isError = false
		await OrderService.create()
			.then(data => (id = data))
			.catch(error => {
				isError = true
				if (!error.response) {
					setError(error.response.data.Error)
				}
			})
		if (!isError) {
			toast.success(`Order with id ${id} was successfully created`, {
				position: toast.POSITION.BOTTOM_LEFT
			})
			setUpdateQuery(!updateQuery)
		}
	}

	const displayOrders = orders?.map(
		order => {
			return (
				<OrderItem
					updateQuery={updateQuery}
					setUpdateQuery={setUpdateQuery}
					order={order}
					key={order.OrderId}
				/>
			)
		}
	)

	return (
		<Layout title='Orders'>
			<table className={styles.table}>
				<thead>
				<tr>
					<th className={styles.textLeft}>Order Id</th>
					<th className={styles.textLeft}>Date</th>
					<th className={styles.textLeft}>Price</th>
					<th className={styles.textLeft}>Status</th>
					<th hidden={!isAdmin} className={styles.textLeft}>User Id</th>
				</tr>
				</thead>
				<tbody className='table-hover'>
				{displayOrders}
				</tbody>
			</table>
			<ToastContainer />
		</Layout>
	)
}

export default Orders
