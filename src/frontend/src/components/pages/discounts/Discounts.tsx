import React, { FC, useState } from 'react'
import Home from '../home/Home'
import { UserService } from '../../../services/user.service'
import { useQuery } from 'react-query'
import { InstrumentService } from '../../../services/instrument.service'
import InstrumentItem from '../../ui/instrument-item/InstrumentItem'
import Layout from '../../ui/layout/Layout'
import styles from '../discounts/Discounts.module.scss'
import ReactPaginate from 'react-paginate'
import Modal from '../../ui/modal/Modal'
import AddInstrumentDBForm from '../../ui/form/AddInstrumentDBForm'
import { DiscountService } from '../../../services/discount.service'
import DiscountItem from '../../ui/discount-item/DiscountItem'
import { ToastContainer } from 'react-toastify'
import AddDiscountDBForm from '../../ui/form/AddDiscountDBForm'

const Discounts: FC = () => {
	let isAdmin =
		UserService.getCurrentIsAdmin() != null &&
		UserService.getCurrentIsAdmin() == 'true'
	const [AddDiscountInDBModalActive, setAddDiscountInDBModalActive] =
		useState(false)

	const { data: discounts, isLoading } = useQuery(
		['discounts'],
		() => DiscountService.getList(),
		{
			select: ({ discounts }) => discounts
		}
	)

	const [pageNumber, setPageNumber] = useState(0)
	const discountsPerPage = 24
	const pagesVisited = pageNumber * discountsPerPage

	const displayDiscounts = discounts
		?.slice(pagesVisited, pagesVisited + discountsPerPage)
		.map(discount => {
			let userId = UserService.getCurrentUserId()
			if (userId == null) return
			if (!isAdmin && discount.UserId != parseInt(userId)) return
			return <DiscountItem discount={discount} key={discount.DiscountId} />
		})

	// @ts-ignore
	const pagesCount = Math.ceil(discounts?.length / discountsPerPage)

	// @ts-ignore
	const changePage = ({ selected }) => {
		setPageNumber(selected)
	}

	return (
		<Layout title='Discounts'>
			<div>
				<button
					className={styles.openBtn}
					onClick={() => setAddDiscountInDBModalActive(true)}
					hidden={!isAdmin}
				>
					Add discount in DB
				</button>
				<table className={styles.table}>
					<thead>
						<tr>
							<th className={styles.textLeft}>Instrument Id</th>
							<th hidden={!isAdmin} className={styles.textLeft}>
								User Id
							</th>
							<th className={styles.textLeft}>Amount</th>
							<th className={styles.textLeft}>Type</th>
							<th className={styles.textLeft}>DateBegin</th>
							<th className={styles.textLeft}>DateEnd</th>
							<th hidden={!isAdmin} className={styles.textLeft}>
								Links
							</th>
						</tr>
					</thead>
					<tbody className='table-hover'>{displayDiscounts}</tbody>
				</table>
				<ReactPaginate
					previousLabel={'Prev'}
					nextLabel={'Next'}
					pageCount={pagesCount}
					onPageChange={changePage}
					containerClassName={styles.paginationBttns}
					previousClassName={'previousBttns'}
					nextLinkClassName={'nextBttn'}
					disabledClassName={'paginationDisable'}
					activeClassName={styles.paginationActive}
				/>
				<Modal
					active={AddDiscountInDBModalActive}
					setActive={setAddDiscountInDBModalActive}
				>
					<AddDiscountDBForm />
				</Modal>
			</div>
			<ToastContainer />
		</Layout>
	)
}

export default Discounts
