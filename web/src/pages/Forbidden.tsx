const ForbiddenPage = () => {
  return (
    <div className="flex items-center justify-center h-screen bg-gray-100">
      <div className="text-center">
        <h1 className="text-6xl font-bold text-yellow-500">403</h1>
        <p className="mt-4 text-xl text-gray-700">Forbidden</p>
        <p className="text-gray-500">
          You don't have permission to access this page.
        </p>
      </div>
    </div>
  );
};

export default ForbiddenPage;
