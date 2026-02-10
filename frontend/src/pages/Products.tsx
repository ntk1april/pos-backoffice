import React, { useState, useEffect } from 'react';
import Layout from '../components/Layout';
import ProductTable from '../components/ProductTable';
import ProductModal from '../components/ProductModal';
import { useAuth } from '../context/AuthContext';
import { productApi } from '../api/products';
import { transactionApi } from '../api/transactions';
import { storeApi, Store } from '../api/stores';
import { Product, CreateProductRequest, UpdateProductRequest } from '../types';

const Products: React.FC = () => {
    const { isAdmin } = useAuth();
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    // Pagination
    const [page, setPage] = useState(1);
    const [pageSize] = useState(10);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(0);

    // Search
    const [search, setSearch] = useState('');
    const [searchInput, setSearchInput] = useState('');

    // Modal
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [modalMode, setModalMode] = useState<'create' | 'edit'>('create');
    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);

    // Stock adjustment
    const [stockModalOpen, setStockModalOpen] = useState(false);
    const [stockProduct, setStockProduct] = useState<Product | null>(null);
    const [stockType, setStockType] = useState<'increase' | 'decrease'>('increase');
    const [stockQuantity, setStockQuantity] = useState(0);
    const [stockUnitPrice, setStockUnitPrice] = useState(0);
    const [stockStoreId, setStockStoreId] = useState<number | null>(null);
    const [stockNotes, setStockNotes] = useState('');
    const [stores, setStores] = useState<Store[]>([]);

    useEffect(() => {
        loadProducts();
        loadStores();
    }, [page, search]);

    const loadStores = async () => {
        try {
            const data = await storeApi.getStores();
            setStores(data);
        } catch (err) {
            console.error('Failed to load stores:', err);
        }
    };

    const loadProducts = async () => {
        setLoading(true);
        setError('');
        try {
            const data = await productApi.getProducts(page, pageSize, search, 'ACTIVE');
            setProducts(data.products);
            setTotal(data.total);
            setTotalPages(data.total_pages);
        } catch (err: any) {
            setError(err.response?.data?.message || 'Failed to load products');
        } finally {
            setLoading(false);
        }
    };

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        setSearch(searchInput);
        setPage(1);
    };

    const handleCreateProduct = () => {
        setModalMode('create');
        setSelectedProduct(null);
        setIsModalOpen(true);
    };

    const handleEditProduct = (product: Product) => {
        setModalMode('edit');
        setSelectedProduct(product);
        setIsModalOpen(true);
    };

    const handleSaveProduct = async (productData: CreateProductRequest | UpdateProductRequest) => {
        if (modalMode === 'create') {
            await productApi.createProduct(productData as CreateProductRequest);
        } else if (selectedProduct) {
            await productApi.updateProduct(selectedProduct.id, productData as UpdateProductRequest);
        }
        loadProducts();
    };

    const handleDeleteProduct = async (id: number) => {
        if (window.confirm('Are you sure you want to delete this product?')) {
            try {
                await productApi.deleteProduct(id);
                loadProducts();
            } catch (err: any) {
                alert(err.response?.data?.message || 'Failed to delete product');
            }
        }
    };

    const handleStockAdjust = (product: Product, type: 'increase' | 'decrease') => {
        setStockProduct(product);
        setStockType(type);
        setStockQuantity(0);
        setStockUnitPrice(type === 'increase' ? product.cost || 0 : product.price || 0);
        setStockStoreId(null);
        setStockNotes('');
        setStockModalOpen(true);
    };

    const handleStockSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!stockProduct || stockQuantity <= 0 || stockUnitPrice <= 0) return;

        // Validate: DECREASE requires store selection
        if (stockType === 'decrease' && !stockStoreId) {
            alert('Please select a store for sales transactions');
            return;
        }

        try {
            await transactionApi.createTransaction({
                transaction_type: stockType === 'increase' ? 'INCREASE' : 'DECREASE',
                product_id: stockProduct.id,
                store_id: stockType === 'decrease' ? stockStoreId : null,
                quantity: stockQuantity,
                unit_price: stockUnitPrice,
                notes: stockNotes,
            });

            setStockModalOpen(false);
            loadProducts();
        } catch (err: any) {
            alert(err.response?.data?.message || 'Failed to adjust stock');
        }
    };

    return (
        <Layout>
            <div className="px-4 py-6 sm:px-0">
                <div className="flex justify-between items-center mb-6">
                    <h1 className="text-3xl font-bold text-gray-900">Products</h1>
                    {isAdmin && (
                        <button
                            onClick={handleCreateProduct}
                            className="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                        >
                            Add Product
                        </button>
                    )}
                </div>

                {/* Search */}
                <form onSubmit={handleSearch} className="mb-6">
                    <div className="flex gap-2">
                        <input
                            type="text"
                            placeholder="Search products..."
                            value={searchInput}
                            onChange={(e) => setSearchInput(e.target.value)}
                            className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                        />
                        <button
                            type="submit"
                            className="px-6 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                        >
                            Search
                        </button>
                        {search && (
                            <button
                                type="button"
                                onClick={() => {
                                    setSearch('');
                                    setSearchInput('');
                                    setPage(1);
                                }}
                                className="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300"
                            >
                                Clear
                            </button>
                        )}
                    </div>
                </form>

                {/* Error */}
                {error && (
                    <div className="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
                        {error}
                    </div>
                )}

                {/* Loading */}
                {loading ? (
                    <div className="text-center py-12">
                        <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
                        <p className="mt-4 text-gray-600">Loading products...</p>
                    </div>
                ) : (
                    <>
                        {/* Table */}
                        <ProductTable
                            products={products}
                            onEdit={handleEditProduct}
                            onDelete={handleDeleteProduct}
                            onStockAdjust={handleStockAdjust}
                        />

                        {/* Pagination */}
                        {totalPages > 1 && (
                            <div className="mt-6 flex items-center justify-between">
                                <div className="text-sm text-gray-700">
                                    Showing page {page} of {totalPages} ({total} total products)
                                </div>
                                <div className="flex gap-2">
                                    <button
                                        onClick={() => setPage(page - 1)}
                                        disabled={page === 1}
                                        className="px-4 py-2 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                                    >
                                        Previous
                                    </button>
                                    <button
                                        onClick={() => setPage(page + 1)}
                                        disabled={page === totalPages}
                                        className="px-4 py-2 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                                    >
                                        Next
                                    </button>
                                </div>
                            </div>
                        )}
                    </>
                )}

                {/* Product Modal */}
                <ProductModal
                    isOpen={isModalOpen}
                    onClose={() => setIsModalOpen(false)}
                    onSave={handleSaveProduct}
                    product={selectedProduct}
                    mode={modalMode}
                />

                {/* Stock Adjustment Modal */}
                {stockModalOpen && (
                    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                        <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
                            <div className="px-6 py-4 border-b border-gray-200">
                                <h2 className="text-2xl font-bold text-gray-900">
                                    {stockType === 'increase' ? 'Buy from Supplier' : 'Sell to Store'}
                                </h2>
                                <p className="text-sm text-gray-600 mt-1">
                                    {stockProduct?.name} (Current Stock: {stockProduct?.stock})
                                </p>
                            </div>

                            <form onSubmit={handleStockSubmit} className="px-6 py-4">
                                {/* Store Selection (only for DECREASE) */}
                                {stockType === 'decrease' && (
                                    <div className="mb-4">
                                        <label className="block text-sm font-medium text-gray-700 mb-1">
                                            Select Store <span className="text-red-500">*</span>
                                        </label>
                                        <select
                                            value={stockStoreId || ''}
                                            onChange={(e) => setStockStoreId(parseInt(e.target.value))}
                                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                                            required
                                        >
                                            <option value="">-- Select a store --</option>
                                            {stores.map(store => (
                                                <option key={store.id} value={store.id}>
                                                    {store.name} ({store.code})
                                                </option>
                                            ))}
                                        </select>
                                    </div>
                                )}

                                <div className="mb-4">
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Quantity <span className="text-red-500">*</span>
                                    </label>
                                    <input
                                        type="number"
                                        min="1"
                                        value={stockQuantity}
                                        onChange={(e) => setStockQuantity(parseInt(e.target.value))}
                                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                                        required
                                    />
                                </div>

                                <div className="mb-4">
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Unit Price (฿) <span className="text-red-500">*</span>
                                    </label>
                                    <input
                                        type="number"
                                        min="0.01"
                                        step="0.01"
                                        value={stockUnitPrice}
                                        onChange={(e) => setStockUnitPrice(parseFloat(e.target.value))}
                                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                                        required
                                    />
                                    <p className="text-xs text-gray-500 mt-1">
                                        {stockType === 'increase' ? 'Cost price from supplier' : 'Selling price to store'}
                                    </p>
                                </div>

                                {/* Total Amount Display */}
                                <div className="mb-4 p-3 bg-gray-50 rounded-md">
                                    <div className="flex justify-between items-center">
                                        <span className="text-sm font-medium text-gray-700">Total Amount:</span>
                                        <span className="text-lg font-bold text-gray-900">
                                            ฿{(stockQuantity * stockUnitPrice).toLocaleString()}
                                        </span>
                                    </div>
                                </div>

                                <div className="mb-4">
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Notes
                                    </label>
                                    <textarea
                                        value={stockNotes}
                                        onChange={(e) => setStockNotes(e.target.value)}
                                        rows={3}
                                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                                        placeholder="Optional notes..."
                                    />
                                </div>

                                <div className="flex justify-end space-x-3">
                                    <button
                                        type="button"
                                        onClick={() => setStockModalOpen(false)}
                                        className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
                                    >
                                        Cancel
                                    </button>
                                    <button
                                        type="submit"
                                        className={`px-4 py-2 text-sm font-medium text-white rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 ${stockType === 'increase'
                                            ? 'bg-green-600 hover:bg-green-700 focus:ring-green-500'
                                            : 'bg-orange-600 hover:bg-orange-700 focus:ring-orange-500'
                                            }`}
                                    >
                                        {stockType === 'increase' ? 'Increase' : 'Decrease'} Stock
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}
            </div>
        </Layout>
    );
};

export default Products;
