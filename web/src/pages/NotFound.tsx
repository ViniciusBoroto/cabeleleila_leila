const NotFoundPage = () => {
  return (
    <div className="flex items-center justify-center h-screen bg-gray-100">
      <div className="text-center">
        <h1 className="text-6xl font-bold text-red-500">404</h1>
        <p className="mt-4 text-xl text-gray-700">Oops! Page not found.</p>
        <p className="text-gray-500">
          The page you are looking for might have been removed, or is
          temporarily unavailable.
        </p>
      </div>
    </div>
  );
};

export default NotFoundPage;
